layui.use('form', function () {
    var form = layui.form;

    form.on('radio(type)', function (data) {
        changeType(data.value)
    });
});

function changeType(t) {
    t += '';
    var title = $("#post-title");
    var desc = $("#post-desc");
    switch (t) {
        case '1':
            title.removeClass("hide");
            desc.removeClass("hide");
            break;
        case '2':
            title.addClass("hide");
            desc.removeClass("hide");
            break;
        case '3':
            title.addClass("hide");
            desc.addClass("hide");
            break;
        default:
            console.error('错误的文章类型');
            break
    }
}

var ce = $('#ContentEditor');
var cp = $('#ContentPreview');
ce.on("input", function () {
    convertMD();
});


function convertmd() {
    // https://prismjs.com/extending.html#highlight-all
    // cp.html(converter.makeHtml(ce.val()));
    // Prism.highlightAll(true);

    $.ajax({
        type: "POST",
        contentType: "application/json",
        url: "/admin/posts/markdown",
        context: document.body,
        data: ce.val(),
        success: function (result) {
            cp.html(result);
            Prism.highlightAll(true);
        },
        error: function (result) {
            console.log(result);
            toast('解析markdown失败');
        }
    });
}

var convertMD = function debounce(fn, delay = 2000) {
    let timer;

    return function () {
        var context = this;
        var args = arguments;

        clearTimeout(timer);

        timer = setTimeout(function () {
            fn.apply(context, args);
        }, delay);
    }
}(convertmd, 1000);

var tags;
var post;

$.ajax({
    type: "GET",
    contentType: "application/json",
    url: "/admin/tags/list",
    context: document.body,
    success: function (result) {
        tags = result.result;
        renderIfBoth();
    },
    error: function (result) {
        console.log(result);
        toast('加载列表失败');
    }
});


var arr = window.location.href.split('/');
id = arr[arr.length - 1];
if (isNaN(id)) {
    $('#postForm').attr('action', '/admin/posts');
    post = {'create_flag': true};
    renderIfBoth();
} else {
    $('#postForm').attr('action', '/admin/posts/' + id);

    $.ajax({
        type: "GET",
        contentType: "application/json",
        url: "/admin/posts/" + id,
        context: document.body,
        success: function (result) {
            console.log(result);
            post = result.result;
            renderIfBoth();
        },
        error: function (result) {
            console.log(result);
            toast('加载文章失败');
        }
    });
}

function renderIfBoth() {
    if (tags && post) {
        render();
    }
}

function render() {
    layui.use(['form', 'laytpl'], function () {
        var laytpl = layui.laytpl;
        var form = layui.form;

        //================= post =============
        if (!post.create_flag) {
            form.val("postForm", {
                "Type": post.Type,
                "Title": post.Title,
                "Description": post.Description,
                "Cover": post.Cover,
                "Content": post.Content,
                "Publish": post.Publish ? '1' : '0',
            });

            changeType(post.Type);

            if (post.TagIDs) {
                var arr = post.TagIDs.split(",");
                tags.forEach(function (item, index) {
                    item.checked = false;
                    for (let s of arr) {
                        if ((item.ID + "") === s) {
                            item.checked = true;
                            break
                        }
                    }
                });
            }
        }

        //================= tag =============
        var tt = document.querySelector('#tags');
        var tpl = `<input type="checkbox" id="tag-{{d.ID}}" value="{{d.ID}}" name="TagIDs" title="{{d.Name}}" {{# if(d.checked){ }}checked{{# } }} >\n`;
        tags.forEach(function (item, index) {
            laytpl(tpl).render(item, function (html) {
                tt.innerHTML += html;
            })
        });

        form.render();
        convertMD();
    })
}