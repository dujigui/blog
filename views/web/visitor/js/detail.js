layui.use('form', function(){
    var form = layui.form;

    //监听提交
    form.on('submit(formDemo)', function(data){
        layer.msg(JSON.stringify(data.field));
        console.log('shit');
        return false;
    });

    console.log('fuck');
});

console.log('cao');