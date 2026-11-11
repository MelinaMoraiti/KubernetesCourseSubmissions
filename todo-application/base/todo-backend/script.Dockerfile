FROM redis:7-alpine

COPY create-todo.sh /usr/local/bin/create-todo.sh

RUN chmod +x /usr/local/bin/create-todo.sh

ENTRYPOINT ["/usr/local/bin/create-todo.sh"]