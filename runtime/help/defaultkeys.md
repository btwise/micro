# 默认键

下面是默认热键及其功能的简单图表。更多的
请提供有关绑定自定义热键或更改默认绑定的信息
运行' > help keybindings '

请记住,这里的*所有*键都是可重定义的!如果你不喜欢,你可以更改它!

### 超级用户

  键                           功能描述                                                                          
---------- |--------------------------------------------------------------------------------------------------
 Ctrl-e    | 打开命令提示符以运行命令(有关有效命令列表,请参阅' >help commands').                                      
 Tab       | 在命令提示符中,如果可能,它将自动完成.                                                                 
 Ctrl-b    | 运行一个shell命令(这将在执行命令时关闭micro).                          

### 导航

   键                          功能描述                                                                                    
---------------------------------------------------------------------------------------------------------------------- 
上下左右键                    | 移动光标                                                                                   
Shift-箭头                   | 移动并选择文本                                                                              
Alt(Mac上是Ctrl)-左箭头       | 移动到当前行的开头                                                                           
Alt(Mac上是Ctrl)-右箭头       | 移动到当前行的结尾                                                                           
Home键                       | 移动到当前行文本的开头                                                                       
End键                        | 移动到当前行的末尾                                                                          
Ctrl(Mac是Alt)-左箭头         | 将光标向左移动一个单词                                                                       
Ctrl(Mac是Alt)-右箭头         | 向右移动光标一个单词                                                                         
Alt-{                        | 移动光标到上一个空行,或文档的开头                                                             
Alt-}                       | 将光标移动到下一个空行或文档末尾                                                               
PageUp                      | 将光标上翻一页                                                                              
PageDown                    | 将光标下翻一页                                                                              
Ctrl-Home or Ctrl-UpArrow   | 移动光标到文档的开头                                                                         
Ctrl-End or Ctrl-DownArrow  | 移动光标到文档的结尾                                                                         
Ctrl-l                      | 跳转到文件中的一行(使用#提示符)                                                                
Ctrl-w                      | 在当前选项卡的分隔窗口之间循环 (使用 `> vsplit` 或 `> hsplit` 创建分割窗口)                       

### Tabs

 键          功能描述   
-------- |-------------------------- 
 Ctrl-t  | 打开一个新标签              
 Alt-,   | 上一个标签                 
 Alt-.   | 下一个标签                 

### 查找操作

 键                   功能描述                         
---------- |------------------------------------------ 
 Ctrl-f    | 查找 (打开提示)                             
 Ctrl-n    | 查找当前搜索的下一个                         
 Ctrl-p    | 查找当前搜索的上一个                         

注意:Ctrl-n和Ctrl-p应该从主缓冲区中使用,而不是从内部使用
搜索提示符。在Ctrl-f后,按enter键完成搜索,然后你可以使用
Ctrl-n和Ctrl-p来循环匹配.

### 文件操作

   键           功能描述                                           
---------- |------------------------------------------------------------------ 
 Ctrl-q    | 关闭当前文件(如果这是最后一个打开的文件,则退出micro)    
 Ctrl-o    | 打开一个文件(提示输入文件名)                                
 Ctrl-s    | 保存当前文件                                                 

### 文本操作

        键                                  功能描述                   
------------------------------------ |------------------------------------------ 
 Ctrl(Alt on Mac)-Shift-RightArrow   | 选择右边的单词                         
 Ctrl(Alt on Mac)-Shift-LeftArrow    | 选择左边的单词                           
 Alt(Ctrl on Mac)-Shift-LeftArrow    | 选择当前行开始           
 Alt(Ctrl on Mac)-Shift-RightArrow   | 选择当前行结束            
 Shift-Home                          | 选择当前行开始           
 Shift-End                           | 选择当前行结束             
 Ctrl-Shift-UpArrow                  | 选择文件的开始                  
 Ctrl-Shift-DownArrow                | 选择文件的结尾                     
 Ctrl-x                              | 剪切选定文本                         
 Ctrl-c                              | 拷贝选定文本                        
 Ctrl-v                              | 粘贴                                     
 Ctrl-k                              | 剪切当前行                          
 Ctrl-d                              | 重复当前行                    
 Ctrl-z                              | 撤销                                      
 Ctrl-y                              | 重做                                      
 Alt-UpArrow                         | 向上移动当前行或选定的行    
 Alt-DownArrow                       | 向下移动当前行或选定的行  
 Alt-Backspace or Alt-Ctrl-h         | 删除左边单词                          
 Ctrl-a                              | 全选                                
 Tab                                 | 缩进所选文本                      
 Shift-Tab                           | 取消缩进所选文本                    

### 宏

     键              功能描述                                                           
---------- |---------------------------------------------------------------------------------- 
 Ctrl-u    | 切换宏录制(按Ctrl-u开始录制,再按一次停止录制)  
 Ctrl-j    | 运行最新录制的宏                                                         

### 多光标

   键                     功能描述                                                                       
------------------ |---------------------------------------------------------------------------------------------- 
 Alt-n             | 从选区中创建新的多个游标(如果没有当前选区,将选择当前单词)  
 Alt-Shift-Up      | 在当前光标上方的行上生成一个新光标                                          
 Alt-Shift-Down    | 在当前光标下面的行上生成一个新光标                                          
 Alt-p             | 删除最新的多个光标                                                                 
 Alt-c             | 删除所有多个游标(取消)                                                          
 Alt-x             | 跳过多个光标选择                                                                
 Alt-m             | 在当前选定的每一行的开始处生成一个新游标                    
 Ctrl-MouseLeft    | 在任意位置放置多个光标                                                       

### 其他

 键                功能描述                                                               
---------- |-------------------------------------------------------------------------------------- 
 Ctrl-g    | 打开帮助文件                                                                        
 Ctrl-h    | 退格键(旧终端不支持退格键,使用Ctrl+H代替)     
 Ctrl-r    | 切换行号标尺                                                          

### Emacs风格的操作

 键             功能描述   
---------- |-------------------------- 
 Alt-f     | 下一个单词                 
 Alt-b     | 上一个单词             
 Alt-a     | 移动到行开始     
 Alt-e     | 移动到行结尾       

### 功能键

警告!功能键可能无法在所有终端使用! 

 键            功能描述   
------ |-------------------------- 
 F1    | 打开帮助                 
 F2    | 保存                      
 F3    | 查找                      
 F4    | 退出                      
 F7    | 查找                      
 F10   | 退出                      
