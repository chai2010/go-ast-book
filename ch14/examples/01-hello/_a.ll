@printf_format_int = constant [4 x i8] c"%d\0A\00"

declare i32 @printf(i8* %format, ...)

define i32 @main() {
entry:
	%0 = getelementptr [4 x i8], [4 x i8]* @printf_format_int, i32 0, i32 0
	%1 = call i32 (i8*, ...) @printf(i8* %0, i32 42)
	ret i32 0
}
