# /backend/convert_writings_to_pdf.py (优化版)

import os
from pathlib import Path
import markdown
from weasyprint import HTML, CSS
import comtypes.client  # 引入 comtypes
import time

# --- 配置 ---
BACKEND_DIR = Path(__file__).resolve().parent
WRITINGS_DIR = BACKEND_DIR / "static" / "writings"
WD_FORMAT_PDF = 17  # Word 中 PDF 格式的常量

def convert_md_to_pdf(md_path, pdf_path):
    """使用 markdown 和 WeasyPrint 将 Markdown 文件转换为 PDF (此函数不变)"""
    print(f"  [MD] 正在转换: {md_path.name} -> {pdf_path.name}")
    try:
        md_content = md_path.read_text(encoding='utf-8')
        html_content = markdown.markdown(md_content, extensions=['fenced_code', 'tables'])
        css_style = """
            @page { size: A4; margin: 2cm; }
            body { font-family: sans-serif; line-height: 1.6; }
            h1, h2, h3 { color: #333; border-bottom: 1px solid #eee; padding-bottom: 5px; }
            pre { background-color: #f4f4f4; padding: 1em; border-radius: 5px; overflow-x: auto; }
            code { font-family: monospace; }
            table { border-collapse: collapse; width: 100%; }
            th, td { border: 1px solid #ddd; padding: 8px; }
            th { background-color: #f2f2f2; }
            img { max-width: 100%; height: auto; }
        """
        html = HTML(string=html_content)
        css = CSS(string=css_style)
        html.write_pdf(pdf_path, stylesheets=[css])
        print("  ✅ 转换成功！")
        return True
    except Exception as e:
        print(f"  ❌ 转换失败: {e}")
        return False

def batch_convert_docx_to_pdf(docx_files):
    """
    批量转换 DOCX 文件。只启动一次 Word 实例来处理所有文件。
    """
    word = None
    successful_conversions = []
    
    try:
        print("\n▶️  正在启动 Microsoft Word 实例...")
        # 启动 Word 应用程序，并设置为不可见
        word = comtypes.client.CreateObject("Word.Application")
        word.Visible = False
        print("✅ Word 实例已在后台启动。")

        for docx_path, pdf_path in docx_files:
            print(f"  [DOCX] 正在处理: {docx_path.name}")
            doc = None
            try:
                # 打开文档
                doc = word.Documents.Open(str(docx_path.resolve()))
                # 保存为 PDF
                doc.SaveAs(str(pdf_path.resolve()), FileFormat=WD_FORMAT_PDF)
                print(f"    -> 成功保存为 {pdf_path.name}")
                successful_conversions.append(docx_path)
            except Exception as e:
                print(f"    ❌ 处理失败: {e}")
            finally:
                # 确保关闭当前文档，即使转换失败
                if doc:
                    doc.Close()
                # 短暂延时，给 Word 一点喘息时间
                time.sleep(0.5)

    except Exception as e:
        print(f"❌ 启动或操作 Word 实例时发生严重错误: {e}")
        print("   请确保 Microsoft Word 已正确安装。")
    finally:
        # 确保最终关闭 Word 应用程序
        if word:
            print("\n▶️  正在关闭 Microsoft Word 实例...")
            word.Quit()
            print("✅ Word 实例已关闭。")
    
    return successful_conversions

def main():
    print(f"▶️  开始扫描文章目录: {WRITINGS_DIR.resolve()}")
    if not WRITINGS_DIR.exists():
        print(f"❌ 错误: 目录不存在 {WRITINGS_DIR}")
        return

    converted_count = 0
    updated_count = 0
    skipped_count = 0
    
    docx_to_process = []
    
    for file_path in WRITINGS_DIR.iterdir():
        if not file_path.is_file():
            continue

        ext = file_path.suffix.lower()
        
        # 检查是否需要转换
        should_process = False
        status = ""
        
        if ext in [".docx", ".md"]:
            pdf_path = file_path.with_suffix(".pdf")
            if not pdf_path.exists():
                print(f"发现新文件: {file_path.name}")
                should_process = True
                status = "new"
            elif file_path.stat().st_mtime > pdf_path.stat().st_mtime:
                print(f"发现更新文件: {file_path.name}")
                should_process = True
                status = "update"
            else:
                skipped_count += 1
                continue
        
        if should_process:
            if ext == ".md":
                if convert_md_to_pdf(file_path, pdf_path):
                    if status == "new": converted_count += 1
                    else: updated_count += 1
            elif ext == ".docx":
                # 先收集所有需要处理的 DOCX 文件
                docx_to_process.append((file_path, pdf_path, status))

    # --- 统一处理所有 DOCX 文件 ---
    if docx_to_process:
        print(f"\n发现 {len(docx_to_process)} 个 DOCX 文件需要处理。")
        # 提取路径信息进行批量转换
        docx_paths_to_convert = [(item[0], item[1]) for item in docx_to_process]
        successfully_converted_paths = batch_convert_docx_to_pdf(docx_paths_to_convert)
        
        # 根据成功列表更新计数
        for docx_path, _, status in docx_to_process:
            if docx_path in successfully_converted_paths:
                if status == "new": converted_count += 1
                else: updated_count += 1

    print("\n--- 扫描完成 ---")
    print(f"✅ 新增转换 {converted_count} 个文件。")
    print(f"🔄 更新了 {updated_count} 个文件。")
    print(f"⏭️ 跳过了 {skipped_count} 个已是最新版本的文件。")

if __name__ == "__main__":
    main()