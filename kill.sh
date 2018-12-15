(ps aux | grep ./build/auth-serv | awk 'NR==1{print $2}' >&3 ) 3>pid
kill -9 $(<pid)
(ps aux | grep ./build/game-serv | awk 'NR==1{print $2}' >&3 ) 3>pid
kill -9 $(<pid)
(ps aux | grep ./build/api-serv | awk 'NR==1{print $2}' >&3 ) 3>pid
kill -9 $(<pid)
