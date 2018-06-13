FROM alpine:latest
LABEL secure alpine python dev env

# Install openrc
RUN apk update && apk add openrc &&\
# Tell openrc its running inside a container, till now that has meant LXC
    sed -i 's/#rc_sys=""/rc_sys="lxc"/g' /etc/rc.conf &&\
# Tell openrc loopback and net are already there, since docker handles the networking
    echo 'rc_provide="loopback net"' >> /etc/rc.conf &&\
# no need for loggers
    sed -i 's/^#\(rc_logger="YES"\)$/\1/' /etc/rc.conf &&\
# can't get ttys unless you run the container in privileged mode
    sed -i '/tty/d' /etc/inittab &&\
# can't set hostname since docker sets it
    sed -i 's/hostname $opts/# hostname $opts/g' /etc/init.d/hostname &&\
# can't mount tmpfs since not privileged
    sed -i 's/mount -t tmpfs/# mount -t tmpfs/g' /lib/rc/sh/init.sh &&\
# can't do cgroups
    sed -i 's/cgroup_add_service /# cgroup_add_service /g' /lib/rc/sh/openrc-run.sh

# make sure the package repository is up to date
RUN apk update && apk upgrade

# add bash
RUN apk add bash
RUN sed -i -e "s/bin\/ash/bin\/bash/" /etc/passwd

# add the other basics such as vim, sshd and python2.7
RUN apk add vim openssh python sntpc sudo py-pip
RUN pip install --upgrade pip
RUN pip install virtualenv

# configure sntpc time to sync once per day
RUN > /etc/conf.d/sntpc
RUN printf '%s\n' '# Set time sync interval (in seconds) here.' 'SNTPC_INTERVAL="86400"' 'SNTPC_HOST=pool.ntp.org' >/etc/conf.d/sntpc

# set the sshd_config
RUN echo "PubkeyAcceptedKeyTypes=+ssh-dss" > /etc/ssh/ssh_config
RUN echo "PasswordAuthentication yes" > /etc/ssh/sshd_config \
        && echo "Protocol 2" >> /etc/ssh/sshd_config\
        && echo "ChallengeResponseAuthentication no" >> /etc/ssh/sshd_config

# add the services to startup with openrc
RUN rc-update add sshd
RUN rc-update add sntpc

# setup the dev user with no password but a public key
ENV USER dev
RUN adduser -D -g "" -s /bin/bash $USER
USER dev
ADD authorized_keys /home/$USER/.ssh/authorized_keys

USER root
RUN echo "dev:12345" | chpasswd
RUN echo "dev ALL=(ALL) ALL" >> /etc/sudoers

EXPOSE 22
ENV LC_ALL=en_US.UTF-8

CMD ["/sbin/init"]
