---
date modified: Sunday, September 10th 2023, 8:59:05 am
tags:
  - daily-note
type: daily
liquid: 0
sleep: 0
steps: 0
alcohol: false
beef: false
chicken: false
coffein: false
milk: false
date: "{{ .Current.DateOnly }}"
diarrhea: false
fart: false
fruit: false
headache: false
month: "[[{{ .DailyNoteFolder }}/{{ .Current.Year }}/{{ .Current.Month }}/00 Index|{{ .Current.Month }}]]"
significant: false
vegetables: false
year: "[[{{ .DailyNoteFolder }}/{{ .Current.Year }}/00 Index|{{ .Current.Year }}]]"
watch charged: false
work time: 0
work location: {{ .WorkLocation }}
weight: 0
---

<< [[{{ .Previous.DateOnly }}|Previous]] | [[00 Index|Index]] | [[{{ .Next.DateOnly }}|Next]] >>

# Daily Log

## Table of contents

- [[{{ .Current.DateOnly }}#Today's meetings|Today's meetings]]
- [[{{ .Current.DateOnly }}#Today's birthdays|Today's birthdays]]
- [[{{ .Current.DateOnly }}#Work|Work]]
- [[{{ .Current.DateOnly }}#Health|Health]]
  - [[{{ .Current.DateOnly }}#Food & beverages|Food & beverages]]
    - [[{{ .Current.DateOnly }}#Beverages|Beverages]]
    - [[{{ .Current.DateOnly }}#Food|Food]]
- [[{{ .Current.DateOnly }}#Other stuff|Other stuff]]
- [[{{ .Current.DateOnly }}#Tasks|Tasks]]
- [[{{ .Current.DateOnly }}#New Items Created|New Items created]]

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
where birthday and dateformat(date(birthday), "MM-dd") = "{{ .Current.Month }}-{{ .Current.Day }}"
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
done "{{ .Current.DateOnly }}"
```

## New or changed items

```dataview
table date-created as "Planted at",
date-modified as "Last tended to",
length(file.inlinks) as "In Links", 
length(file.outlinks) as "Out Links"
where (date(file.cday) <= (date(this.date) + dur(1 day))
and date(file.cday) >= date(this.date)) or ((date(file.mday) <= (date(this.date) + dur(1 day))
and date(file.mday) >= date(this.date)))
and contains(file.name, "Calls") = false
and contains(file.name, "Messages") = false
```
