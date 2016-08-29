FROM alpine
ADD raft /home/raft
ENTRYPOINT ["/home/raft"]
