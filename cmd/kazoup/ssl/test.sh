#!/bin/bash
micro --enable_tls --tls_cert_file=cert.pem --tls_key_file=key.pem --registry=mdns web
