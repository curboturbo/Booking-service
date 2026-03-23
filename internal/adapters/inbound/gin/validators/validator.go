package validators

func CheckDaysOfWeek(days []int) bool {
    n := len(days)
    if n == 0 || n > 7 {
        return false
    }
    for i := 0; i < n; i++ {
        if days[i] < 1 || days[i] > 7 {
            return false
        }
        if i > 0 && days[i] <= days[i-1] {
            return false
        }
    }
    return true
}