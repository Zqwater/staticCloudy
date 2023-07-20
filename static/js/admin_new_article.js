  function submitForm(event) {
    event.preventDefault();
    const form = document.querySelector('form');
    const formData = new FormData(form);
    
    const requestOptions = {
      method: 'POST',
      body: formData
    };

    fetch('/admin/articles', requestOptions)
      .then(response => {
        if (!response.ok) {
          throw new Error('请求格式不正确');
        }
        return response.json();
      })
      .then(data => {
        console.log(data);
        window.location.href = '/admin/articles';
      })
      .catch(error => {
        console.error(error);
      });
  }

  function previewImage(event) {
    const fileInput = event.target;
    const imagePreview = document.getElementById('preview');
  
    if (fileInput.files && fileInput.files[0]) {
      const reader = new FileReader();
  
      reader.onload = function(e) {
        imagePreview.src = e.target.result;
        imagePreview.style.display = 'block';
      }
  
      reader.readAsDataURL(fileInput.files[0]);
    } else {
      imagePreview.src = '';
      imagePreview.style.display = 'none';
    }
  }