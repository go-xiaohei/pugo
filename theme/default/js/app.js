$(document).ready(function () {

    $('pre code').each(function (i, block) {
        $(block).addClass("line-numbers");
        Prism.highlightElement(block, false)
    });

    $(".scroll-to").click(function (event) {
        event.preventDefault();
        console.log("scroll to:", $(event.target).attr("href"));
        //get the top offset of the target anchor
        var target_offset = $($(event.target).attr("href")).offset();
        var target_top = target_offset.top;
        $('html, body').animate({scrollTop: target_top}, 500, 'easeInSine');
    });

    $('#comment-user').val(localStorage.getItem("comment-user"));
    $('#comment-email').val(localStorage.getItem("comment-email"));
    $('#comment-url').val(localStorage.getItem("comment-url"));
    $('#comment-form').ajaxForm({
        dataType: "json",
        success: function (json) {
            if (!json.status) {
                $('#comment-error').text(json.error).show();
                setTimeout(function () {
                    $('#comment-error').hide()
                }, 2000);
                return;
            }
            var html = $('#comment-template').html();
            html = html.replace("{avatar_url}", json.data.comment.avatar, -1);
            html = html.replace("{name}", json.data.comment.name, -1);
            html = html.replace("{name}", json.data.comment.name, -1);
            html = html.replace("{url}", json.data.comment.url, -1);
            html = html.replace("{content}", json.data.comment.body, -1);
            html = html.replace("{created}", json.data.comment.created, -1);
            if (json.data.comment.status != 1) {
                html = html.replace("{status}", "disapproved red-text", -1);
                html = html.replace("{status}", "disapproved", -1);
            }
            $('.comment-header').after(html);
            $('#comment-content').val("");
            localStorage.setItem("comment-user", $('#comment-user').val());
            localStorage.setItem("comment-email", $('#comment-email').val());
            localStorage.setItem("comment-url", $('#comment-url').val());
        }
    });

    $('.comment-section').on("click", ".reply", function () {
        var $comment = $($(this).attr("href"));
        var name = $comment.find("> .body > .meta .author").text();
        var body = $comment.find(".content").html();
        $("#comment-reply-view").html('<div class="name">@' + name + '</div><blockquote>' + body + '</blockquote>').show();
        $('#comment-parent').val($comment.data("id"));
        // scroll to form
        var top = $('#comment-form').offset().top;
        $('body,html').animate({scrollTop: top}, 250);
        // show cancel button
        $('#comment-reply-btn').show();
    });
    $('#comment-reply-btn').on("click", function () {
        $('#comment-parent').val(0);
        $("#comment-reply-view").html("").hide();
        $(this).hide();
    });

    $(".comment").each(function (i, item) {
        var $item = $(item);
        var p = $item.data("parent");
        if (p > 0) {
            var $pComment = $('#comment-' + p);
            if ($pComment.length) {
                $pComment.find(".children").append($item);
            } else {
                $item.remove();
            }
        }
    });
});
