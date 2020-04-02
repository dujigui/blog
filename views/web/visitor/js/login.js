layui.use(['form'], function () {
    $ = layui.jquery
    $('#qq').click(function() {
        var url = `https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=${qqAppId}&redirect_uri=${encodeURI(redirectUri)}&state=${qqState}`;
        console.log(url);
        window.open(url);
    })
});



function google() {

}