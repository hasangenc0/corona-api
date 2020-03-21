# docker build -t corona .
# docker run -it --rm -p 8080:8080 corona
heroku container:login
heroku create
heroku container:push web
heroku container:release web