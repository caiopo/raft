FROM ubuntu
ADD test_raft /home/test_raft
CMD ["/home/test_raft"]