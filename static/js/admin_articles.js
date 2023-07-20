  function confirmDelete(articleID) {
    if (confirm("确定要删除这篇文章吗？")) {
      deleteArticle(articleID);
    }
  }
  
  function deleteArticle(articleID) {
    fetch(`/admin/articles/${articleID}`, {
      method: 'DELETE'
    })
    .then(response => {
      if (response.ok) {
        return response.json();
      }
      throw new Error('无法删除文章');
    })
    .then(data => {
      alert(data.message); // Show the server response message
      if (data.success) {
        window.location.reload();
      }
    })
    .catch(error => {
      console.error(error);
      alert('删除文章失败，请稍后重试。');
    });
  }