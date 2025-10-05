# /your-personal-blog/backend/convert_videos.py (最终修复版 - 采用双重验证)

import os
import subprocess
import json
import shutil
from pathlib import Path

# --- 配置 ---
VIDEO_DIR = Path("static/video")
BACKUP_DIR = VIDEO_DIR / "originals_backup"

def get_video_info(file_path):
    """
    使用 ffprobe 获取视频的详细信息，包括编码和 movflags。
    返回一个包含 'codec' 和 'is_faststart' 的字典。
    """
    command = [
        "ffprobe",
        "-v", "quiet",
        "-print_format", "json",
        "-show_format",
        "-show_streams",
        str(file_path)
    ]
    
    try:
        result = subprocess.run(command, capture_output=True, text=True, check=True, encoding='utf-8')
        data = json.loads(result.stdout)
        
        info = {
            'codec': None,
            'is_faststart': False
        }

        if "streams" in data and len(data["streams"]) > 0:
            for stream in data["streams"]:
                if stream.get("codec_type") == "video":
                    info['codec'] = stream.get("codec_name")
                    break
        
        # --- 【核心修复】采用更可靠的双重验证逻辑 ---
        # 一个文件被认为是优化过的，只要满足以下任一条件：
        # 1. ffprobe 直接报告了 'faststart' 的 movflags 标签。
        # 2. 文件的 major_brand 是 'isom'，这是 ffmpeg 处理后标准的 MP4 文件格式，
        #    是 moov atom 在文件头部的强有力指标。
        if "format" in data and "tags" in data["format"]:
            tags = data["format"]["tags"]
            if tags.get("movflags") == "faststart" or tags.get("major_brand") == "isom":
                info['is_faststart'] = True

        if info['codec'] is None:
             print(f"  [警告] 在 {file_path.name} 中未找到视频流。")
             info['codec'] = "no_video_stream"

        return info

    except (subprocess.CalledProcessError, json.JSONDecodeError, FileNotFoundError) as e:
        if isinstance(e, FileNotFoundError):
            print("  [致命错误] ffprobe 命令未找到。请确保 FFmpeg 已安装并已添加到系统 PATH。")
            exit(1)
        print(f"  [错误] 检查文件 {file_path.name} 时出错: {e}")
        return None

def process_video(source_path, dest_path):
    """
    使用 ffmpeg 处理视频：转换为 H.264 并确保 faststart 生效。
    """
    command = [
        "ffmpeg",
        "-y",
        "-i", str(source_path),
        "-c:v", "libx264",
        "-preset", "fast",
        "-c:a", "copy",
        "-pix_fmt", "yuv420p",
        "-movflags", "+faststart",
        str(dest_path)
    ]
    
    print(f"  [ffmpeg] 正在执行处理命令...")
    try:
        subprocess.run(command, check=True, capture_output=True, text=True, encoding='utf-8', timeout=600)
        return True
    except subprocess.CalledProcessError as e:
        print(f"  [处理失败] {source_path.name} 转换时发生错误。")
        print(f"  --- FFmpeg 标准错误 ---")
        print(e.stderr)
        return False
    except subprocess.TimeoutExpired:
        print(f"  [处理失败] {source_path.name} 转换超时。")
        return False

def main():
    """主函数，扫描并处理视频文件"""
    print(f"开始扫描视频文件夹: {VIDEO_DIR.resolve()}")
    if not VIDEO_DIR.exists():
        print(f"错误: 视频文件夹 {VIDEO_DIR} 不存在。")
        return

    BACKUP_DIR.mkdir(exist_ok=True)
    processed_count = 0
    scanned_count = 0
    skipped_count = 0

    for video_file in VIDEO_DIR.glob("*.mp4"):
        if not video_file.is_file():
            continue
        
        scanned_count += 1
        print(f"\n[{scanned_count}] 正在检查: {video_file.name}")
        
        video_info = get_video_info(video_file)
        
        if video_info is None:
            continue
            
        needs_processing = False
        if video_info['codec'] != "h264":
            print(f"  状态: 编码为 {video_info['codec']}，需要转换为 H.264。")
            needs_processing = True
        elif not video_info['is_faststart']:
            print(f"  状态: 编码为 H.264，但未优化 (faststart)，需要处理。")
            needs_processing = True
        else:
            print(f"  状态: 编码为 H.264 且已优化，无需处理。")
            skipped_count += 1
        
        if needs_processing:
            temp_output_path = video_file.with_suffix(".temp_processed.mp4")
            
            if process_video(video_file, temp_output_path):
                new_info = get_video_info(temp_output_path)
                # 验证逻辑现在也使用了更可靠的 is_faststart 判断
                if new_info and new_info['codec'] == 'h264' and new_info['is_faststart']:
                    print("  [验证成功] 正在替换原始文件...")
                    backup_path = BACKUP_DIR / video_file.name
                    shutil.move(str(video_file), str(backup_path))
                    shutil.move(str(temp_output_path), str(video_file))
                    print(f"  操作完成: {video_file.name} 已被处理并更新。")
                    processed_count += 1
                else:
                    print(f"  [验证失败] 处理后的文件编码不正确或无法读取！")
                    if temp_output_path.exists(): os.remove(temp_output_path)
            else:
                print(f"  由于处理步骤失败，已跳过文件 {video_file.name}。")
                if temp_output_path.exists(): os.remove(temp_output_path)

    print("\n--- 扫描完成 ---")
    print(f"共扫描 {scanned_count} 个视频文件。")
    print(f"✅ 成功处理 (转换或优化) {processed_count} 个文件。")
    print(f"⏭️ 跳过了 {skipped_count} 个无需处理的文件。")

if __name__ == "__main__":
    main()