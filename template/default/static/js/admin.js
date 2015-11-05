$(function () {

    // login form action
    (function () {
        $('#login-form').ajaxForm({
            success: function (json) {
                if (json.status) {
                    window.location.href = "/";
                    return;
                }
                $('#login-error').text(json.error).show();
            }
        })
    })();

    // simplemde editor
    (function () {
        console.log($('.simplemde').length);
        if (!window.SimpleMDE || $('.simplemde').length < 1) {
            return;
        }
        var simplemde = new SimpleMDE({
            element: $('.simplemde')[0],
            spellChecker: false,
            status: false,
            autoDownloadFontAwesome: false,
            autosave: {
                enabled: false,
                unique_id: "article-md",
                delay: 1000
            },
            renderingConfig: {
                singleLineBreaks: false,
                codeSyntaxHighlighting: true
            },
            tabSize: 4,
            initialValue: $('#article-content,#page-content').val()
        });
    })();

    // all ajax form action
    (function () {
        var ajaxFormAction = function (prefix) {
            $('#' + prefix + '-write').ajaxForm({
                beforeSubmit: function () {
                    $('#' + prefix + '-write-success').hide();
                    $('#' + prefix + '-write-error').hide();
                },
                success: function (json) {
                    if (json.status) {
                        $('#' + prefix + '-write-success').show();
                        if(json.data.id) {
                            $('#' + prefix + '-id').val(json.data.id);
                        }
                        return;
                    }
                    $('#' + prefix + '-write-error').text(json.error).show();
                }
            })
        };
        ajaxFormAction("article");
        ajaxFormAction("page");
        ajaxFormAction("general");
    })();

    // article remove button click
    (function () {
        $('.article .remove-btn,#admin-articles .tab-pane .remove-btn').on("click", function () {
            return confirm("Are you sure to delete the article/page ?")
        })
    })();
});
