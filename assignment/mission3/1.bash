sshpass -p 'Hustee108108.' ssh unique-ali 'cd hzy; rm m'
sshpass -p 'Hustee108108.' sftp unique-ali << EOF
put m hzy/
exit
EOF
sshpass -p 'Hustee108108.' ssh unique-ali