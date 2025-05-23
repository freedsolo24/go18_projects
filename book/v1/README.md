# book api
1. 构建gin引擎实对, 并且运行
2. gin引擎实对, 调用方法实现5个api, 先不实现里面的调用函数
   * List book
   ```sh
   GET /api/books
   ```
   * Create book
   ```sh
   POST /api/books
   ```
   * Get book by book number
   ```sh
   GET /api/books/:bn
   ```
   * Update book
   ```sh
   PUT /api/books/:bn
   ```
   * Delete book
   ```sh
   DELETE /api/books/:bn
   ```
3. 在5个api接口里面获取用户的请求参数, 处理, 持久化, 响应.
   * List book使用GET方法, 请求参数在(1)querystring
   * Create book使用POST方法, 请求参数在(1)body
   * List book by number使用GET方法, 请求参数在(1)uri
   * Update book使用PUT方法, 请求参数(1)uri (2)body
   * Delete book使用Delete方法, 请求参数(1)uri
4. 改造成bookapihandler定义以上5个接口
5. 定义book结构体, 映射打tag 映射gorm

