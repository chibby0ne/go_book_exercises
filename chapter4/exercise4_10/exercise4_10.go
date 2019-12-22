// Modify issues to report the rsult in age categories, say less than a month old, less than a year old and more than a year old.

package main

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter4/github"
	"log"
	"os"
	"time"
)

func printIssues(issues []*github.Issue) {
	for _, issue := range issues {
		fmt.Printf("#%-15s #%-5d %9.9s %.55s\n", issue.CreatedAt.Format(time.ANSIC), issue.Number, issue.User.Login, issue.Title)
	}
}

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()
	fmt.Printf("%d issues:\n", result.TotalCount)

	var issuesLessOneMonthOld []*github.Issue
	var issuesMoreOneMonthOld []*github.Issue
	var issuesMoreOneYearOld []*github.Issue

	for _, item := range result.Items {
		if item.CreatedAt.Before(now.AddDate(-1, 0, 0)) {
			issuesMoreOneYearOld = append(issuesMoreOneYearOld, item)
		} else if item.CreatedAt.After(now.AddDate(0, -1, 0)) {
			issuesLessOneMonthOld = append(issuesLessOneMonthOld, item)
		} else {
			issuesMoreOneMonthOld = append(issuesMoreOneMonthOld, item)
		}
	}

	if length := len(issuesLessOneMonthOld); length > 0 {
		fmt.Printf("%d issues < 1 month old:\n", length)
		printIssues(issuesLessOneMonthOld)
	}

	if length := len(issuesMoreOneMonthOld); length > 0 {
		fmt.Printf("%d issues > 1 month old:\n", length)
		printIssues(issuesMoreOneMonthOld)
	}

	if length := len(issuesMoreOneYearOld); length > 0 {
		fmt.Printf("%d issues > 1 year old:\n", length)
		printIssues(issuesMoreOneYearOld)
	}
}
