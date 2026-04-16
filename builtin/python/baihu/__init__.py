import os
from .notify import notify as _notify

def notify(title, text):
    """
    发送内建通知。
    会在调用时校验环境变量：BHPKG_NOTIFY_TOKEN, BHPKG_NOTIFY_CHANNEL
    """
    _TOKEN = os.environ.get("BHPKG_NOTIFY_TOKEN")
    _CHANNEL = os.environ.get("BHPKG_NOTIFY_CHANNEL")

    if not _TOKEN or not _CHANNEL:
        missing = []
        if not _TOKEN: missing.append("BHPKG_NOTIFY_TOKEN")
        if not _CHANNEL: missing.append("BHPKG_NOTIFY_CHANNEL")
        
        error_msg = f"缺少必要的环境变量以使用 baihu 模块: {', '.join(missing)}。请在白虎面板的任务设置中配置指定的 Key。"
        raise RuntimeError(error_msg)
    
    return _notify(title, text)

__all__ = ['notify']
