import os
import uuid
import time

def rename_videos_in_directory(directory_path):
    """
    将指定目录下的所有文件重命名为随机且唯一的名称（保留原始扩展名）。
    增加了安全确认步骤。

    参数:
    directory_path (str): 目标目录的路径。
    """
    try:
        # 检查目录是否存在
        if not os.path.isdir(directory_path):
            print(f"错误：目录不存在！")
            print(f"脚本尝试寻找路径: '{directory_path}'")
            print("请确保脚本位于 'backend' 文件夹内，并且 'static/video' 文件夹存在。")
            return

        # 获取目录下所有文件的列表
        files_to_rename = [f for f in os.listdir(directory_path) if os.path.isfile(os.path.join(directory_path, f))]

        if not files_to_rename:
            print(f"目录 '{directory_path}' 中没有找到任何文件。")
            return

        # --- 安全确认步骤 ---
        print("--------------------------------------------------")
        print(f"脚本将在以下目录中重命名文件：")
        print(f"'{directory_path}'\n")
        print("将要被重命名的文件列表：")
        for filename in files_to_rename:
            print(f"- {filename}")
        print("--------------------------------------------------")
        
        # 请求用户确认
        confirm = input("以上文件将被重命名为随机名称，是否继续？ (输入 'y' 或 'yes' 继续): ").lower()

        if confirm not in ['y', 'yes']:
            print("\n操作已取消。")
            return
        
        print("\n开始重命名...")
        time.sleep(1)

        # 遍历并重命名文件
        count = 0
        for filename in files_to_rename:
            old_file_path = os.path.join(directory_path, filename)
            
            file_root, file_extension = os.path.splitext(filename)
            new_filename = f"{uuid.uuid4()}{file_extension}"
            new_file_path = os.path.join(directory_path, new_filename)

            os.rename(old_file_path, new_file_path)
            print(f"成功: '{filename}' -> '{new_filename}'")
            count += 1

        print(f"\n处理完成！总共有 {count} 个文件被成功重命名。")

    except PermissionError:
        print("\n错误：权限不足。请检查是否有权限修改该文件夹中的文件，或尝试以管理员身份运行此脚本。")
    except Exception as e:
        print(f"\n处理过程中发生未知错误: {e}")

# --- 主程序入口 ---
if __name__ == "__main__":
    # --- 动态路径构建 ---
    # 获取脚本文件自身所在的目录的绝对路径
    # __file__ 是一个特殊变量，代表当前脚本的文件名
    script_directory = os.path.dirname(os.path.abspath(__file__))
    
    # 基于脚本的目录，构建到 'static/video' 文件夹的相对路径
    # os.path.join 会自动处理不同操作系统下的路径分隔符（'/' 或 '\'）
    target_directory = os.path.join(script_directory, "static", "video")
    
    # 运行主函数
    rename_videos_in_directory(target_directory)

    input("\n按 Enter 键退出...")