FROM fedora:29

ADD paymentapp /usr/bin/

ENTRYPOINT ["/usr/bin/paymentapp"]

CMD ["start"]