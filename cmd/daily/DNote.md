---
date modified: Sunday, September 10th 2023, 8:59:05 am
tags:
  - daily-note
type: daily
systole:
  - 0
diastole:
  - 0
systole diff: 0
diastole diff: 0
liquid: 0
sleep: 0
steps: 0
alcohol: false
beef: false
chicken: false
coffein: false
date: "{{date:YYYY-MM-DD}}"
diarrhea: false
fruit: false
headache: false
month: "[[00-09 Personal/00 Daily Notes/{{date:YYYY}}/{{date:MM}}/00 Index|{{date:MM}}]]"
significant: false
vegetables: false
year: "[[00-09 Personal/00 Daily Notes/{{date:YYYY}}/00 Index|{{date:YYYY}}]]"
watch charged: false
work time: 0
work location: Home
weight: 0
---

<% tp.user.generate_periodic_link_markdown(tp, "D") %>

# Daily Log

## Table of contents

- [[{{date:YYYY-MM-DD}}#Today's meetings|Today's meetings]]
- [[{{date:YYYY-MM-DD}}#Today's birthdays|Today's birthdays]]
- [[{{date:YYYY-MM-DD}}#Work|Work]]
- [[{{date:YYYY-MM-DD}}#Health|Health]]
	- [[{{date:YYYY-MM-DD}}#Food & beverages|Food & beverages]]
		- [[{{date:YYYY-MM-DD}}#Beverages|Beverages]]
		- [[{{date:YYYY-MM-DD}}#Food|Food]]
- [[{{date:YYYY-MM-DD}}#Other stuff|Other stuff]]
- [[{{date:YYYY-MM-DD}}#Tasks|Tasks]]
- [[{{date:YYYY-MM-DD}}#New Items Created|New Items created]]

## Today's meetings

```dataview
TABLE
dateformat(date(date), "HH:mm") as Time,
title as Title
FROM #meeting
WHERE dateformat(date(date), "yyyy-MM-dd") = dateformat(date(this.date), "yyyy-MM-dd")
SORT dateformat(date(date), "HH:mm") ASC
```

# Today's birthdays

```dataview
TABLE
birthday
FROM #contact
where birthday and dateformat(date(birthday), "MM-dd") = "{{date:MM-DD}}"
```

## Work

- This contains a list of stuff I did related to work

## Health

Wake up:

### Food & beverages

#### Beverages

#### Food

## Other stuff

- This contains a list of stuff I did not related to work

## Tasks

```tasks
group by function task.tags.length?task.tags.map((tag)=>(tag=="#task/todo"?"Me":(tag=="#task/work"?"Work":(tag=="#task/DMREG"?"Work":(tag=="#task/comm"?"Me":(tag=="#task/it"?"Work":(tag=="#task/bug"?"Dev":(tag=="#task/stream"?"Watch":(tag=="#task/shopping"?"Buy":(tag=="#task/carmen"?"Carmen":(tag=="#task/dev"?"Dev":"Me"))))))))))):"Me"
hide tags
sort by urgency
not done
path does not include Template
path does not include group_task
starts before tomorrow
```

### Tasks done

```tasks
group by function task.tags.length?task.tags.map((tag)=>(tag=="#task/todo"?"Me":(tag=="#task/work"?"Work":(tag=="#task/DMREG"?"Work":(tag=="#task/comm"?"Me":(tag=="#task/it"?"Work":(tag=="#task/bug"?"Dev":(tag=="#task/stream"?"Watched":(tag=="task/shopping"?"Buy":(tag=="#task/carmen"?"Carmen":(tag=="#task/dev"?"Dev":"Me"))))))))))):"Me"
hide tags
sort by urgency
done "{{date:YYYY-MM-DD}}"
```

## New or changed items

```dataview
table date-created as "Planted at",
date-modified as "Last tended to",
length(file.inlinks) as "In Links", 
length(file.outlinks) as "Out Links"
where (date(file.cday) <= (date(this.file.cday) + dur(1 day))
and date(file.cday) >= date(this.file.cday)) or (date(file.mday) <= (date(this.file.mday) + dur(1 day))
and date(file.mday) >= date(this.file.mday))
and contains(file.name, "Calls") = false
and contains(file.name, "Messages") = false
```
