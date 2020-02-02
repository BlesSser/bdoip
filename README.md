#代码结构
>server
>>msgChannel  
>>NetworkAdapter和Repository之间交互的通道
>
>>networkAdapter  
>>处理网络请求：  
>>1. 接收request报文，合并Message，放入inputChannel  
>>2. 读取outputChannel，拆分Message，返回
>>>listener  
>>>1. 监听端口，获取数据包  
>>>2. 向client发送Response数据包给client
>> 
>>>assembler  
>>>1. 合并request Message，放入inputChannel  
>>>2. 读取outputChannel，按需拆分给listener返回
>
>>repository
>>1. 读取inputChannel的Message，处理
>>2. 生成Response Message，放入outputChannel
>>
>>>certificate  
>>>负责消息的加解密、签名、验证  
>>>
>>>coder  
>>>Message Byte和Message struct的编解码 
>>>  
>>>operator  
>>>消息的具体处理模块
>>>

