layui.use(['element'], function () {
    $ = layui.jquery;

    $( ".search-box-editor" ).on("input", function() {
        var keyword = $(this).val().toLowerCase();
        $( ".tags" ).each(function () {
            tag = $(this);
            if (!keyword) {
                tag.removeClass("hide");
                tag.find('li').each(function () {
                    $(this).removeClass("hide");
                });
            } else if(tag.attr('id').toLowerCase().includes(keyword)) {
                tag.removeClass("hide");
                tag.find('li').each(function () {
                    $(this).removeClass("hide");
                });
            } else {
                var contained = false;
                tag.find('a').each(function () {
                    a = $(this);
                    if(a.text().toLowerCase().includes(keyword)) {
                        contained = true;
                        a.parent().removeClass("hide");
                    } else {
                        a.parent().addClass("hide");
                    }
                });
                if (contained) {
                    tag.removeClass("hide");
                } else {
                    tag.addClass("hide");
                }
            }
        })
    });
});