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
blood pressure high: false
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

```dynamic-embed
[[Daily appointment]]
```

## Today's birthdays

```dynamic-embed
[[Daily birthday]]
```

## Work

## Health

### Food & beverages

#### Beverages

#### Food

## Volt

## Praxis

## Other stuff

## Tasks done

```tasks
group by function task.tags.length?task.tags.map((tag)=>(tag=="#task/todo"?"Me":(tag=="#task/work"?"Work":(tag=="#task/DMREG"?"Work":(tag=="#task/comm"?"Me":(tag=="#task/it"?"Work":(tag=="#task/bug"?"Dev":(tag=="#task/stream"?"Watched":(tag=="task/shopping"?"Buy":(tag=="#task/carmen"?"Carmen":(tag=="#task/dev"?"Dev":"Me"))))))))))):"Me"
hide tags
sort by urgency
((done "{{ .Current.DateOnly }}") OR (cancelled "{{ .Current.DateOnly }}"))
```

## New or changed items

```dynamic-embed
[[New or changed items]]
```
