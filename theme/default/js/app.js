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

    $('#comment-form').ajaxForm({
        dataType:"json",
        success:function(json){
            if(!json.status){
                $('#comment-error').text(json.error).show();
                setTimeout(function(){
                    $('#comment-error').hide()
                },2000);
                return;
            }
            var html = $('#comment-template').html();
            html = html.replace("{avatar_url}",json.data.comment.avatar,-1);
            html = html.replace("{name}",json.data.comment.name,-1);
            html = html.replace("{name}",json.data.comment.name,-1);
            html = html.replace("{url}",json.data.comment.url,-1);
            html = html.replace("{content}",json.data.comment.body,-1);
            html = html.replace("{created}",json.data.comment.created,-1);
            if(json.data.comment.status != 1){
                html = html.replace("{status}","disapproved red-text",-1);
                html = html.replace("{status}","disapproved",-1);
            }
            $('.comment-header').after(html);
        }
    })
});
