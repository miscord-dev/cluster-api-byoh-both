module.exports = ({
    github,
    context,
    core,
    glob,
    io,
    exec,
    fetch,
    require,
}) => {
    return context.payload.client_payload.value
}
