FROM ubuntu:latest
LABEL maintainer="Liliya233"
RUN sed -i 's/archive.ubuntu.com/mirrors.bfsu.edu.cn/g' /etc/apt/sources.list \
    && apt-get -yq update \
    && apt-get install ca-certificates -y \
    && apt-get install libreadline8 -y \
    && apt-get install gcc -y \
    && apt-get install wget -y \
    && apt-get install -y tzdata \
    && ln -snf /usr/share/zoneinfo/$TIME_ZONE /etc/localtime && echo $TIME_ZONE > /etc/timezone \
    && dpkg-reconfigure -f noninteractive tzdata \
    && mkdir /ome && cd /ome \
#    && wget https://***/res/launcher/omega_launcher_linux_amd64 \
    && wget https://github.com/Liliya233/omega_launcher/releases/latest/download/omega_launcher_linux_amd64 \
    && chmod +x /ome/omega_launcher_linux_amd64 \
    && mkdir /workspace \
    && apt clean -yq && apt autoclean -yq && apt autoremove -yq && rm -rf /var/lib/apt/lists/* \
    && echo "#!/bin/bash" >> /ome/launcher_liliya \
    && echo "sleep 1" >> /ome/launcher_liliya \
    && echo "clear" >> /ome/launcher_liliya \
    && echo 'echo -e "\033[1;33m注意: 如果使用MCSM, 请在终端设置中打开仿真终端，并直接在命令行中输入内容\033[0m"' >> /ome/launcher_liliya \
    && echo "cd /workspace" >> /ome/launcher_liliya \
    && echo "/ome/omega_launcher_linux_amd64" >> /ome/launcher_liliya \
    && chmod +x /ome/launcher_liliya \
    && mkdir -p /root/.config && mkdir -p /root/.config/fastbuilder \
    && echo -n 'zh_CN' > /root/.config/fastbuilder/language
WORKDIR /workspace
ENTRYPOINT [ "/ome/launcher_liliya" ]
