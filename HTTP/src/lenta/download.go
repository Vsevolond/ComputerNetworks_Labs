package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Name, Volume string
}

/*func readItem(item *html.Node) *Item {
	children := getChildren(item)
	for _, child := range children {
		if isElem(child, "a") && getAttr(child, "class") == "n-news_title" {
			return &Item{
				Ref:   getAttr(child, "href"),
				Title: getChildren(child)[0].Data,
			}
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	if isDiv(node, "n-news_list") {
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "n-news_item") {
				if item := readItem(c); item != nil {
					items = append(items, item)
				}
			}
		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}*/

func readItem(item *html.Node) *Item {
	//fmt.Println(os.Stdout, "read item")
	ref := readRef(item)
	name := readName(item)
	volume := readVolume(item)
	fmt.Fprintf(os.Stdout, "%s %s %s\n", ref, name, volume)
	return &Item{
		Ref:    ref,
		Name:   name,
		Volume: volume,
	}
}

func readRef(item *html.Node) string {
	//fmt.Println(os.Stdout, "in ref")
	if getAttr(item, "class") == "cmc-link" {
		return getAttr(item, "href")
	}
	for child := item.FirstChild; child != nil; child = item.NextSibling {
		if ref := readRef(child); ref != "" {
			return ref
		}
	}

	return ""
}

func readName(item *html.Node) string {
	//fmt.Println(os.Stdout, "in name")
	if getAttr(item, "class") == "sc-1eb5slv-0 iworPT" {
		return item.FirstChild.Data
	}
	for child := item.FirstChild; child != nil; child = child.NextSibling {
		if name := readName(child); name != "" {
			return name
		}
	}
	return ""
}

func readVolume(item *html.Node) string {
	//fmt.Println(os.Stdout, "in volume")
	if getAttr(item, "class") == "sc-1eb5slv-0 hykWbK font_weight_500" {
		return item.FirstChild.Data
	}
	for child := item.FirstChild; child != nil; child = child.NextSibling {
		if volume := readVolume(child); volume != "" {
			return volume
		}
	}
	return ""
}

/*func print(items []*Item) {
	for _, item := range items {
		fmt.Println(os.Stdout, item.Name, item.Volume)
	}
}*/

func search(node *html.Node) []*Item {
	if isElem(node, "tbody") {
		//fmt.Println(os.Stdout, "in tbody")
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isElem(c, "tr") {
				//fmt.Println(os.Stdout, "in tr")
				item := readItem(c)
				items = append(items, item)
			}
		}
		//print(items)
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func downloadNews() []*Item {
	log.Println("sending request to coinmarketcap.com")
	if response, err := http.Get("https://coinmarketcap.com"); err != nil {
		log.Println("request to coinmarketcap.com failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Println("got response from coinmarketcap.com", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Println("invalid HTML from coinmarketcap.com", "error", err)
			} else {
				log.Println("HTML from coinmarketcap.com parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
