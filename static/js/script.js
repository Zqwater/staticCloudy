// script.js
function callAPI(apiUrl) {
    fetch(apiUrl)
      .then(response => response.text())
      .then(data => {
        console.log(data);
        // 处理 API 响应数据
      })
      .catch(error => {
        console.error('发生错误:', error);
      });
  }
  
  document.getElementById('search-button').addEventListener('click', function() {
    var searchInput = document.getElementById('search-input');
    var searchTerm = searchInput.value;
  
    // 在这里执行搜索操作，可以调用 callAPI 函数来获取搜索结果
  
    // 示例：调用 callAPI 函数来搜索
    var apiUrl = 'https://api.example.com/search?term=' + searchTerm;
    callAPI(apiUrl);
  });
  
  $(document).ready(function() {
    function changeBackground() {
      var backgrounds = [
        '/static/image/bg/bg1.jpg',
        '/static/image/bg/bg2.jpg',
        '/static/image/bg/bg3.jpg',
        '/static/image/bg/bg4.jpg',
        '/static/image/bg/bg5.jpg'
      ];
      var randomIndex = Math.floor(Math.random() * backgrounds.length);
      var randomBackground = backgrounds[randomIndex];
      document.getElementById('background').style.backgroundImage = "url('" + randomBackground + "')";
    }
  
    changeBackground();
    setInterval(changeBackground, 60000); // Change background every 1 minute
  
    $('.search-icon').click(function() {
      var query = $('.search-bar input[type="text"]').val();
      if (query.trim() !== '') {
        // Perform search operation
        alert('Searching: ' + query);
      }
    });
  
    $('.search-bar input[type="text"]').keydown(function(event) {
      if (event.keyCode === 13) {
        var query = $('.search-bar input[type="text"]').val();
        if (query.trim() !== '') {
          // Perform search operation
          alert('Searching: ' + query);
        }
      }
    });
  
    $('.api-scroll-container').on('click', '.api-item', function() {
      var port = $(this).data('port');
      var currentURL = window.location.href;
      var newURL = currentURL.replace(/:\d+/, ":" + port);
      window.location.href = newURL;
    });
  
    $('.api-scroll-container').on('scroll', function() {
      var scrollLeft = $(this).scrollLeft();
      if (scrollLeft <= 0) {
        $('.api-scroll-left').hide();
      } else {
        $('.api-scroll-left').show();
      }
  
      var containerWidth = $(this).outerWidth();
      var contentWidth = $(this).find('.api-list').outerWidth();
      if (scrollLeft >= contentWidth - containerWidth) {
        $('.api-scroll-right').hide();
      } else {
        $('.api-scroll-right').show();
      }
    });
  
    $('.api-scroll-left').click(function() {
      var scrollContainer = $('.api-scroll-container');
      scrollContainer.animate({ scrollLeft: '-=300' }, 300);
    });
  
    $('.api-scroll-right').click(function() {
      var scrollContainer = $('.api-scroll-container');
      scrollContainer.animate({ scrollLeft: '+=300' }, 300);
    });
  
    // 新增二维码点击事件
    $('.qr-code-item').click(function() {
      var qrCodeAlt = $(this).find('img').attr('alt');
      alert('二维码点击事件: ' + qrCodeAlt);
    });
  
    // 新增二维码放大功能
    $('.qr-code-item img').click(function() {
      var src = $(this).attr('src');
      var alt = $(this).attr('alt');
      var qrCodeImage = $('<img>').attr('src', src).attr('alt', alt).addClass('qr-code-image');
      var qrCodeModal = $('<div>').addClass('qr-code-modal').append(qrCodeImage);
      $('body').append(qrCodeModal);
  
      $('.qr-code-modal').click(function() {
        $(this).remove();
      });
    });
  });
  