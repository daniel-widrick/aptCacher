# Simple Cache for apt/etc

Step 1) run the docker container

step 2) Run the below command replacing localhost with your domain/ip

    cat <<EOF > /etc/apt/apt.conf.d/01proxy
    Acquire::http::Proxy "http://localhost:8700/";
    EOF


