# Thực hành tốt nhất phát triển Golang

Dưới đây là một số thực hành tốt nhất khi phát triển Golang:

## Thư viện

- Sử dụng Shard db cho MySQL, tham khảo Vitess
- Đối với OAuth trong Go, xem xét sử dụng thư viện ory/fosite hoặc ory/hydrax
- Xác thực rất quan trọng trong microservices, do đó nên sử dụng thư viện cho nó
- Truy cập trang web [Effective Go](https://golang.org/doc/effective_go) để biết thêm chi tiết

## Đồng bộ

- Golang là ngôn ngữ đồng bộ, có nghĩa là nó chạy từ trên xuống dưới
- Goroutines nhẹ hơn các luồng và một luồng có thể chạy nhiều goroutines
- Data racing là vấn đề phổ biến khi nhiều goroutines liên tục ghi vào một biến, dẫn đến kết quả không chính xác
- Sử dụng Mutex Lock để xử lý vấn đề này, cụ thể là sync.RWMutex
- Hoặc sử dụng Channels để tránh sử dụng Mutex Lock (được khuyến nghị)
- Mỗi yêu cầu được xử lý trong một goroutine

## Cơ sở dữ liệu

- Để tăng tốc độ tìm kiếm trong cơ sở dữ liệu có quan hệ Many2Many, hãy đánh index bảng có nhiều record hơn trước
- Để tăng tốc độ load page nhất là với infinity scroll thì nên dùng Seek Method (tìm kiếm `<` `>` với id) thay vì OFFSET trong DB
- Preload trong gorm không phải là Join nên sẽ có thể về null nếu table được Preload k có key trên table chính

## Yêu cầu HTTP

- Sử dụng PUT để cập nhật toàn bộ bản ghi và PATCH để chỉ cập nhật một phần
- Trong xử lý lỗi, hãy đưa ra mã và thông điệp dễ hiểu
- Sử dụng mã tùy chỉnh để cho phép thông điệp lỗi phản hồi bằng nhiều ngôn ngữ
- Sử dụng trace_id để theo dõi lỗi
- Xử lý sự cố như là một lỗi thông thường
- Ẩn các lỗi nhạy cảm, nhưng giữ lại nhật ký để gỡ lỗi

## Quản lý bộ nhớ

- Sử dụng con trỏ cho các cấu trúc lớn khi truyền chúng vào hàm để giảm việc sử dụng bộ nhớ và tăng tốc độ thực thi mã
- Ưu tiên các phân bổ nhỏ, ngắn hạn
- Sử dụng sync.Pool cho các đối tượng được sử dụng lại thường xuyên để giảm tải trên trình thu gom rác
- Hãy chú ý đến số lượng goroutines được tạo và hạn chế chúng khi có thể
- Sử dụng kênh đệm (buffered channels) để giảm tranh chấp giữa các goroutines, đặc biệt là khi thao tác kênh mất nhiều thời gian

## Tối ưu hóa

- Tối ưu hóa các hoạt động bị giới hạn bởi CPU bằng cách sử dụng các công cụ lý lịch như pprof để xác định và tối ưu hóa các hoạt động bị giới hạn bởi CPU với thuật toán hoặc cấu trúc dữ liệu hiệu quả
- Giảm thiểu việc khóa (locking) vì nó có thể gây ra tranh chấp và ảnh hưởng đến hiệu năng. Hãy xem xét sử dụng cấu trúc dữ liệu không khóa, các hoạt động nguyên tử hoặc các mẫu đồng bộ khác để giảm tranh chấp khóa
- Tối ưu hóa mẫu truy cập bộ nhớ để cải thiện hiệu suất bằng cách truy cập bộ nhớ theo cách tuyến tính hoặc sử dụng lại bộ nhớ đã truy cập gần đây
- Sử dụng các hàm tích hợp sẵn khi có thể vì chúng thường được tối ưu hóa về hiệu năng

## Kiểm tra và đánh giá hiệu năng

- Thường xuyên kiểm tra và đánh giá hiệu năng mã của bạn để xác định các nút cổ chai hiệu năng và tối ưu hóa chúng bằng cách sử dụng các công cụ như pprof, benchstat và các tính năng đánh giá hiệu năng của gói testing

## Golang không có try/catch nhưng có cơ chế Panic & Recover

### Panic & Recover in GoLang

- Panic is a built-in function that stops the ordinary flow of control and begins panicking
- When function F calls panic, execution of F stops, but deffered functions in F will still executing
- Recover is a built-in function that regains control of panicking
- Recover only useful inside deferred functions

### Cơ chế panic và recover sẽ hoạt động trên 1 current stack trace

- Nên tạo 1 hàm recover chung để handle trường hợp này. Check AppRecover function trong code

```
func willBeErrFunc() {

    go func() { // go routine sẽ tạo ra một stack mới
        defer common.AppRecover() // NOTE: Sử dụng 1 hàm defer để recover panic trong go routine

        // NOTE: panic này sẽ làm program exit luôn vì nó không panicking trên current stack trace có recover()
        panic("Panic not on the current stack trace")
    }()


    panic("Panic on the current stack trace")
    recover()
}

```

## Các điều khác

- Sử dụng container.list để lấy danh sách nâng cao
- Sử dụng sqlx Golang khi hiệu năng là ưu tiên hàng đầu
- Xem xét sử dụng Collection List
- Sử dụng WaitGroup
- **Embed Struct**: https://gobyexample.com/struct-embedding
- Nếu user off cookies ở browser thì session còn chạy được không ?: Ko, vì cookies lưu sessionId, nếu k có Id thì k sử dụng đc session
- Nếu muốn lưu jwt của user vô trong db thì nên lưu phần thứ 3 thôi. Chính là phần signature, nếu revolke hay logout thì chỉ cần xoá record đó thôi

===============================================================================================================================

# Go: Một ngôn ngữ lập trình đơn giản, hiệu quả và dễ dàng sử dụng

Go là một ngôn ngữ lập trình được biên dịch, kiểu dữ liệu tĩnh, được thiết kế với mục đích đơn giản, hiệu quả và dễ dàng lập trình. Go thường được mô tả như một ngôn ngữ lập trình thủ tục với một số tính năng lập trình hướng đối tượng và đồng thời.

## Đặc điểm của Go

1. **Thủ tục**: Go tuân theo mô hình lập trình thủ tục, nơi mà mã được tổ chức thành các hàm hoặc thủ tục thực hiện các nhiệm vụ cụ thể. Các hàm này có thể được gọi để thao tác dữ liệu và thực hiện các nhiệm vụ theo từng bước. Cách tiếp cận này tương tự như các ngôn ngữ như C và Pascal.

2. **Tính năng hướng đối tượng**: Mặc dù Go không có lớp (class), kế thừa, hoặc các cấu trúc OOP truyền thống khác, nhưng nó có cấu trúc dữ liệu (structs), phương thức (methods) và giao diện (interfaces), cho phép bạn viết mã theo phong cách giống như OOP. Go tập trung vào sự kết hợp (composition) thay vì kế thừa và ưu tiên sử dụng giao diện để đạt được đa hình (polymorphism).

3. **Đồng thời (Concurrency)**: Go có hỗ trợ đồng thời (concurrent programming) tích hợp với goroutines, kênh (channels) và câu lệnh `select`. Goroutines là các luồng nhẹ do runtime của Go quản lý, giúp viết mã đồng thời dễ dàng hơn mà không phải đối mặt với những phức tạp của đồng bộ hóa luồng. Kênh cung cấp một cách an toàn để giao tiếp giữa các goroutine, và câu lệnh `select` cho phép bạn xử lý nhiều hoạt động giao tiếp đồng thời.

4. **Quản lý bộ nhớ**: Go có bộ thu gom rác (garbage collection), có nghĩa là runtime của ngôn ngữ sẽ tự động quản lý việc cấp phát và thu hồi bộ nhớ cho bạn. Tính năng này đơn giản hóa việc quản lý bộ nhớ so với các ngôn ngữ như C và C++, nơi bạn cần phải quản lý bộ nhớ thủ công.

5. **Kiểu dữ liệu tĩnh**: Go là một ngôn ngữ kiểu dữ liệu tĩnh, có nghĩa là kiểu dữ liệu của các biến được kiểm tra tại thời điểm biên dịch. Điều này giúp phát hiện lỗi sớm và cải thiện độ tin cậy của mã.

6. **Biên dịch**: Mã Go được biên dịch thành mã máy gốc, thường dẫn đến hiệu suất tốt hơn so với các ngôn ngữ thông dịch. Trình biên dịch Go được thiết kế để nhanh chóng, làm cho quá trình xây dựng nhanh chóng và hiệu quả.

Tóm lại, Go có thể được mô tả như một ngôn ngữ lập trình thủ tục với các tính năng hướng đối tượng và hỗ trợ mạnh mẽ cho lập trình đồng thời. Sự tập trung vào sự đơn giản, hiệu quả và dễ sử dụng của Go khiến nó trở thành lựa chọn phổ biến cho việc phát triển phần mềm hiện đại.
