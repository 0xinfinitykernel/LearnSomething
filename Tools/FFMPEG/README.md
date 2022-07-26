在PC上在MP4视频上添加语音M4A的操作程序

❶ 首先访问ffmpeg官方网站。

❷ 使用ffmpeg,在MP4视频文件中添加音频文件M4A,执行以下命令:
command:ffmpeg -i video.mp4 -i audio.M4A -acodec copy -vcodec copy output.mp4
或者
ffmpeg -i video.mp4 -i audio.M4A -c:v copy -c:a aac -map 0:v:0 -map 1:a:0 output.mp4

❸ 积分可以通过指定-c:v copy来防止手边的MP4文件重新编码。通过指定-map 0:v:0 -map 1:a:0,原来的MP4视频即使文件里包含M4A语音曲目,也能确保使用语音文件的曲目。如果MP4文件不包含语音M4A,则不需要指定-map 0:v:0 -map 1:a:0。