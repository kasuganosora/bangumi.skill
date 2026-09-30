package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.Flags().Bool("today", false, "只显示今天播出的动画")
}

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "查看本周放送表",
	Long: `获取本周每天正在播出的动画列表，按星期一到星期日排列。

不需要令牌即可使用。

使用场景:
  - 查看今天有哪些动画更新
  - 追番用户查看本周追番日程

代理和超时读取 config / --proxy，与其他命令一致。

示例:
  bangumi calendar                # 本周放送表
  bangumi calendar --today        # 只看今天
  bangumi calendar --output json  # JSON 格式输出`,

	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewPublicClient()
		if err != nil {
			return err
		}
		items, err := client.GetCalendar(BackgroundCtx())
		if err != nil {
			return err
		}
		today, _ := cmd.Flags().GetBool("today")
		if today {
			items = filterCalendar(items, bangumiWeekdayID(time.Now()))
			if len(items) == 0 || len(items[0].Items) == 0 {
				msg := "今天没有放送"
				return PrintOutput(map[string]string{"message": msg}, stringerFunc(func() string { return msg }))
			}
		}
		return PrintOutput(items, formatCalendar(items))
	},
}

// bangumiWeekdayID 把时间转成 Bangumi 放送表的星期编号（周一=1，周日=7）。
func bangumiWeekdayID(t time.Time) int {
	wd := int(t.Weekday())
	if wd == int(time.Sunday) {
		return 7
	}
	return wd
}

func filterCalendar(items []api.CalendarItem, weekdayID int) []api.CalendarItem {
	var out []api.CalendarItem
	for _, item := range items {
		if item.Weekday.ID == weekdayID {
			out = append(out, item)
		}
	}
	return out
}

func formatCalendar(items []api.CalendarItem) fmt.Stringer {
	return stringerFunc(func() string {
		var b strings.Builder
		b.WriteString("📅 本周放送表\n")
		b.WriteString(strings.Repeat("─", 60) + "\n")
		for _, item := range items {
			if len(item.Items) == 0 {
				continue
			}
			fmt.Fprintf(&b, "\n【%s (%s)】\n", item.Weekday.CN, item.Weekday.JA)
			for i, s := range item.Items {
				fmt.Fprintf(&b, "  %d. %s", i+1, s.Name)
				if s.NameCN != "" {
					fmt.Fprintf(&b, " (%s)", s.NameCN)
				}
				fmt.Fprintf(&b, "  [ID:%d]", s.ID)
				if s.Rating.Score > 0 {
					fmt.Fprintf(&b, "  ⭐%.1f", s.Rating.Score)
				}
				b.WriteString("\n")
			}
		}
		return b.String()
	})
}
