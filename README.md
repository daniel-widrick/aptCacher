# Simple Cache for apt/etc

Step 1) run the docker container

step 2) Run the below command replacing localhost with your domain/ip

    cat <<EOF > /etc/apt/apt.conf.d/01proxy
    Acquire::http::Proxy "http://localhost:8700/";
    EOF

---

I created this as a drop dead simple caching proxy to locally cache apt packages for times when I am frequently rebuilding docker containers.

I wouldn't trust it in production and the port should not be exposed on untrusted networks. It's basically an open proxy with no authentication.

