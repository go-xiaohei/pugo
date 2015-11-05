$(function () {

    (function () {
        if ($("pre code").length) {
            $('head').append($('<link rel="stylesheet" type="text/css" />').attr('href', theme_link+"/static/css/prism.css"));
            $.getScript(theme_link+"/static/js/prism.min.js",function(){
                $('pre code').each(function (i, block) {
                    $(block).addClass("line-numbers");
                    Prism.highlightElement(block, false)
                });
            });
        }
    })();
});