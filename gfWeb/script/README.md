可以传递 2-3个参数
./count_player_log.sh <file_path> <splite_str>  <Condition>
如下：
./count_player_log.sh /opt/t3/trunk/server/log/game/2021_9_2/service_player_log.log / logid,8/type,0
<file_path>： 要统计的文件对象路径
<splite_str> ： 多个条件要分割的字符串 /
<Condition> 条件字符串 ，用于过滤条件 如条件一：过滤15点的则是 ^15 条件二: 等于 logid,8 条件三： 不等于 type,0 则 -V type,0
            得到输入项 ^15/logid,8/-V type,0/playerid,14603}