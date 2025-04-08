FROM gcr.io/dataflow-templates-base/go-template-launcher-base

ARG WORKDIR=/dataflow/template
RUN mkdir -p ${WORKDIR}

COPY wordcount ${WORKDIR}/wordcount

ENV FLEX_TEMPLATE_GO_BINARY="${WORKDIR}/wordcount"

ENTRYPOINT ["/opt/google/dataflow/go_template_launcher"]
