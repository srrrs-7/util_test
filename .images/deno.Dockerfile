FROM denoland/deno

RUN apt update && apt install unzip
RUN deno add @luca/flag 
