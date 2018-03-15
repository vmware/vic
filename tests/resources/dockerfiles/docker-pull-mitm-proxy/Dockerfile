FROM python:3
COPY archive.tar /
COPY docker-mitm.py /
COPY mitmdump /
RUN chmod +x /mitmdump
EXPOSE 8080
CMD /mitmdump -s docker-mitm.py
