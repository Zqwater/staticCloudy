document.getElementById("loginForm").addEventListener("submit", function(e) {
    e.preventDefault(); // 阻止表单提交的默认行为
  
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
  
    // 将参数拼接为表单格式的字符串
    var formData = new URLSearchParams();
    formData.append("username", username);
    formData.append("password", password);
  
    // 发送表单请求
    fetch("/admin/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      },
      body: formData
    })
      .then(response => response.json())
      .then(data => {
        if (data.success) {
          window.location.href = "/admin/dashboard";
        } else {
          var errorElement = document.createElement("p");
          errorElement.classList.add("error-message");
          errorElement.textContent = data.message;
          document.getElementById("loginForm").appendChild(errorElement);
        }
      })
      .catch(error => {
        console.error("Error:", error);
      });
});
