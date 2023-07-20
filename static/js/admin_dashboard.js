function logout() {
    // 发送 POST 请求以注销登录
    fetch("/admin/logout", {
      method: "POST",
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.success) {
          window.location.href = "/admin";
        } else {
          alert("退出登录失败");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("退出登录失败");
      });
  }
  
  function searchArticles(event) {
    event.preventDefault();
    const searchInput = document.getElementById("search-input").value;
    if (searchInput.trim() !== "") {
      window.location.href = "/admin/articles?search=" + encodeURIComponent(searchInput);
    }
  }
  