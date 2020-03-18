sshpass -p 'Hustee108108.' ssh unique-ali rm hzy/m
sshpass -p 'Hustee108108.' sftp unique-ali << EOF
put m hzy/
exit
EOF
sshpass -p 'Hustee108108.' ssh unique-ali