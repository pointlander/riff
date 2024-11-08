// Copyright 2024 The Riff Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	// FlagFile is the input file
	FlagFile = flag.String("file", "", "the input file")
)

func main() {
	flag.Parse()

	if *FlagFile == "" {
		panic("input required")
	}
	input, err := os.ReadFile(*FlagFile)
	if err != nil {
		panic(err)
	}

	key := os.Getenv("KEY")
	fmt.Println(key)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro")
	resp, err := model.GenerateContent(ctx, genai.Text("Improve the following code:"+string(input)))
	if err != nil {
		panic(err)
	}
	for i, candidate := range resp.Candidates {
		out, err := os.Create(fmt.Sprintf("candidate%d.txt", i))
		for _, part := range candidate.Content.Parts {
			fmt.Fprintf(out, "%s", part)
		}
		err = out.Close()
		if err != nil {
			panic(err)
		}
	}

}
