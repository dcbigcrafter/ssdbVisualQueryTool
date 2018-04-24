/**
 * 工具包
 */

/**
 * 全选
 * 
 * allCkb 全选复选框的id items 复选框的name
 */
$(function() {
    $("#checkAll").click(function() {
        // alert(this.checked);
        if ($(this).is(':checked')) {
            $('input[name="chkBox"]').each(function() {
                // 此处如果用attr，会出现第三次失效的情况
                $(this).prop("checked", true);
            });
        } else {
            $('input[name="chkBox"]').each(function() {
                $(this).removeAttr("checked", false);
            });
            // $(this).removeAttr("checked");
        }
    });
});
