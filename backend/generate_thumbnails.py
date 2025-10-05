# /backend/generate_thumbnails.py (已修复路径问题)

import os
import subprocess
from pathlib import Path

# --- 核心修改：动态构建正确的路径 ---
# 获取脚本文件自身所在的目录 (即 backend 文件夹)
SCRIPT_DIR = Path(__file__).resolve().parent
# 基于脚本目录，构建到 static/video 的路径
VIDEO_DIR = SCRIPT_DIR / "static" / "video"
# --- 修改结束 ---


def generate_thumbnail(video_path, thumbnail_path):
    """使用 ffmpeg 为视频生成一张缩略图"""
    print(f"  > 正在为 {video_path.name} 生成缩略图...")
    command = [
        "ffmpeg",
        "-i", str(video_path),
        "-vframes", "1",
        "-q:v", "2",
        "-f", "image2",
        "-y",
        str(thumbnail_path)
    ]
    
    try:
        subprocess.run(command, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        print(f"  ✅ 成功生成: {thumbnail_path.name}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"  ❌ 生成失败: {video_path.name}")
        print(e.stderr.decode())
        return False
    except FileNotFoundError:
        print("  [致命错误] ffmpeg 命令未找到。请确保 FFmpeg 已安装并已添加到系统 PATH。")
        exit(1)

def main():
    """主函数，扫描视频并生成缺失的缩略图"""
    # 现在这里会打印出正确的绝对路径
    print(f"开始扫描视频文件夹: {VIDEO_DIR.resolve()}")
    
    if not VIDEO_DIR.exists():
        print(f"错误: 视频文件夹 {VIDEO_DIR} 不存在。请检查路径是否正确。")
        return

    generated_count = 0
    scanned_count = 0

    for video_file in VIDEO_DIR.glob("*.mp4"):
        if not video_file.is_file():
            continue
        
        scanned_count += 1
        print(f"\n[{scanned_count}] 正在检查: {video_file.name}")
        
        thumbnail_file = video_file.with_suffix(".jpg")
        
        if thumbnail_file.exists():
            print("  - 状态: 缩略图已存在，跳过。")
        else:
            print("  - 状态: 缩略图不存在，开始生成。")
            if generate_thumbnail(video_file, thumbnail_file):
                generated_count += 1

    print("\n--- 扫描完成 ---")
    print(f"共扫描 {scanned_count} 个视频文件。")
    print(f"新生成了 {generated_count} 张缩略图。")

if __name__ == "__main__":
    main()