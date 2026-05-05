package cmd

import (
	"fmt"
	"strings"

	"github.com/kasuganosora/bangumi.skill/cli/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(calendarCmd)
}

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "查看本周放送表",
	Long: `获取本周每天正在播出的动画列表，按星期一到星期日排列。

不需要令牌即可使用。

使用场景:
  - 查看今天有哪些动画更新
  - 追番用户查看本周追番日程

示例:
  bangumi calendar              # 人性化的本周放送表
  bangumi calendar --json       # JSON 格式输出`,

	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}
		items, err := client.GetCalendar(BackgroundCtx())
		if err != nil {
			return err
		}
		return PrintOutput(items, formatCalendar(items))
	},
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
