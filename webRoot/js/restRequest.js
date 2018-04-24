/**
 * 请求rest服务的ajax
 */

$(function() {
	var key = $("#key"); // key对应的文本框
	var value = $("#value"); // value对应的文本框
	var number = $("#number"); // 查询数量对应的文本框
	var ssdbip = $("#ssdbIP"); // ssdb地址对应的文本框
	var ssdbport = $("#ssdbPort"); // ssdb端口对应的文本框
	var dbIp // 接收cookie中的ssdb地址
	var dbPort // 接收cookie中的ssdb端口
	var dbIpStr = "ssdbIP="; // cookie中ssdb地址的项
	var dbPortStr = "ssdbPort="; // cookie中ssdb端口的项
	// 使用cookie存储数据库连接信息
	var ca = document.cookie.split(';');
	// 读取cookie中存储的ssdb连接信息
	for (var i = 0; i < ca.length; i++) {
		var c = ca[i].trim();
		if (c.indexOf(dbIpStr) == 0) {
			dbIp = c.substring(dbIpStr.length, c.length);
		} else if (c.indexOf(dbPortStr) == 0) {
			dbPort = c.substring(dbPortStr.length, c.length);
		}
	}
	// 判断cookie中是否存在ssdb连接信息
	if ((typeof(dbIp) == "undefined") || (typeof(dbPort) == "undefined")) {
		alert("请先连接数据库");
	} else {
		//将连接输入框的ip端口改为cookie里的
		ssdbip.val(dbIp);
		ssdbport.val(dbPort);
		// 通过cookie中存的地址端口连接ssdb
		$.ajax({
			type: "POST",
			// 访问WebService使用Post方式请求
			contentType: "application/json",
			// WebService 会返回Json类型
			dataType: 'json',
			url: "/rest",
			// 调用WebService的地址和方法名称组合 ---- WsURL/方法名
			data: JSON.stringify({
				com: "connect",
				data: {
					dbIP: dbIp,
					dbPort: dbPort,
				}
			}),
			success: function(result) { // 回调函数，result，返回值
				// console.log(result);
				$("#addWarnMsg").text("");
				if (result.code == "200000") {
					$("#warnSsdbIP").text(dbIp);
					$("#warnSsdbPort").text(dbPort);
					$("#warnMsg").text("数据库已连接！");
					$("#warnSsdbIP").css("color", "green");
					$("#warnSsdbPort").css("color", "green");
					$("#warnMsg").css("color", "green");
				} else {
					$("#warnMsg").text(result.msg);
					$("#warnMsg").css("color", "red");
				}
			}
		});
	}

	// 查询或者插入新数据的ajax响应
	$("#butsubmit").click(function() {
		$.ajax({
			type: "POST",
			// 访问WebService使用Post方式请求
			contentType: "application/json",
			// WebService 会返回Json类型
			dataType: 'json',
			url: "/rest",
			// 调用WebService的地址和方法名称组合 ---- WsURL/方法名
			data: JSON.stringify({
				com: "query",
				data: {
					key: key.val(),
					value: value.val(),
					number: number.val(),
					dbIP: dbIp,
					dbPort: dbPort,
				}
			}),
			success: function(result) { // 回调函数，result，返回值
				var count;
				if (result.code == "200000" && result.data != null) {
					$("#tbody").remove();
					$("#showData").append("<tbody id='tbody'>");
					$.each(result.data,
						function(i, row) {
							count = i + 1;
							$("#showData").append("<tr><td>" + (i + 1) + "</td><td><input name='chkBox' value='" + row.key + "' type='checkbox'/></td><td>" + row.key + "</td><td>" + row.value + "</td></tr>");
						});
					$("#showData").append("</tbody>");
					$("#resultText").text("查询结果：" + count + "条记录");
					$("#addWarnMsg").text("");

					//$("#warnSsdbIP").text(result.data.info.dbIP);
					//$("#warnSsdbPort").text(result.data.info.dbPort);
					//$("#warnMsg").text("数据库已连接！");
					//$("#warnSsdbIP").css("color", "green");
					//$("#warnSsdbPort").css("color", "green");
					//$("#warnMsg").css("color", "green");
				} else if (result.code == "300000") {
					$("#addWarnMsg").text("增加数据成功");
					//$("#warnSsdbIP").text(result.data.info.dbIP);
					//$("#warnSsdbPort").text(result.data.info.dbPort);
					//$("#warnMsg").text("数据库已连接！");
					//$("#warnSsdbIP").css("color", "green");
					//$("#warnSsdbPort").css("color", "green");
					//$("#warnMsg").css("color", "green");

				} else {
					alert(result.msg);
				}
			}
		});
		$("#tbody").remove();
	});

	// 改变数据库连接的响应方法
	$("#changeCon").click(function() {
		$.ajax({
			type: "POST",
			// 访问WebService使用Post方式请求
			contentType: "application/json",
			// WebService 会返回Json类型
			dataType: 'json',
			url: "/rest",
			// 调用WebService的地址和方法名称组合 ---- WsURL/方法名
			data: JSON.stringify({
				com: "connect",
				data: {
					dbIP: ssdbip.val(),
					dbPort: ssdbport.val(),
				}
			}),
			success: function(result) { // 回调函数，result，返回值
				console.log(result);
				$("#addWarnMsg").text("");
				if (result.code == "200000") {
					$("#warnSsdbIP").text(ssdbip.val());
					$("#warnSsdbPort").text(ssdbport.val());
					$("#warnMsg").text("数据库已连接！");
					$("#warnSsdbIP").css("color", "green");
					$("#warnSsdbPort").css("color", "green");
					$("#warnMsg").css("color", "green");
					// 连接成功 新的连接地址 存入cookie
					document.cookie = "ssdbIP=" + ssdbip.val();
					document.cookie = "ssdbPort=" + ssdbport.val();
					dbIp = ssdbip.val();
					dbPort = ssdbport.val();
				} else {
					$("#warnMsg").text(result.msg);
					$("#warnMsg").css("color", "red");
				}
			}
		});
	});

	// 删除所选数据的响应方法
	$("#delete").click(function() {
			var id_array = new Array();
			id_array.push("multi_del");
			$('input[name="chkBox"]:checked').each(function() {
				id_array.push($(this).val()); // 向数组中添加元素
			});
			if (id_array == "multi_del") {
				alert("请勾选所要删除的数据！");
				return
			}
			var idstr = id_array.join(','); // 将数组元素连接起来以构建一个字符串
			var msg = "如果以后再也见不到你，那么祝你早安、午安、晚安。"
			if (confirm(msg)) {
				$.ajax({
					type: "POST",
					// 访问WebService使用Post方式请求
					contentType: "application/json",
					// WebService 会返回Json类型
					dataType: 'json',
					url: "/rest",
					// 调用WebService的地址和方法名称组合 ---- WsURL/方法名
					data: JSON.stringify({
						com: "delete",
						data: {
							key: idstr,
							dbIP: dbIp,
					        dbPort: dbPort,
						}
					}),
					success: function(result) { // 回调函数，result，返回值
						// console.log(result);
						alert(result.msg);
						if (result.code == "200000") {
							$("input[name='chkBox']:checked").each(function() { // 遍历选中的checkbox
								n = $(this).parents("tr").index(); // 获取checkbox所在行的顺序
								$("tbody#tbody").find("tr:eq(" + n + ")").remove();
								$("#checkAll").attr("checked", false);
							});
							chkBoxNum = $('#tbody').children('tr').length;
							$("#resultText").text("查询结果：" + chkBoxNum + "条记录");
						}
					}
				});
			}
	});
})
