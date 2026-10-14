# CineGraph Domain

The shared language for describing narrative knowledge, time, spoilers, and
supporting evidence in CineGraph.

## Language

**Work**:
A film, series, season, episode, or other authored narrative that provides the
scope in which story facts are interpreted.
_Avoid_: Title, content, media

**Entity**:
A distinct subject within or about a Work, such as a character, place,
organization, object, or creative contributor.
_Avoid_: Node, item

**Narrative Fact**:
A proposition about one or more Entities within the scope of a Work.
_Avoid_: Fact, triple, edge

**Narrative Event**:
An occurrence in story time involving one or more Entities.
_Avoid_: Scene, plot point

**Story Time**:
The chronology in which Narrative Events occur inside a Work's world.
_Avoid_: Timeline, event time

**Release Time**:
The chronology in which information became available to an audience through
official publication or release.
_Avoid_: Publication date, air date

**Claim**:
A sourced assertion that a Narrative Fact is true, false, uncertain, or true
only during a stated interval.
_Avoid_: Record, statement

**Evidence**:
A traceable source that supports or challenges a Claim.
_Avoid_: Citation, reference

**Spoiler Horizon**:
The latest point in a Work's Release Time through which a requester permits
knowledge to be revealed.
_Avoid_: Spoiler level, progress

**Spoiler Policy**:
The rules that decide whether a result may be revealed under a Spoiler Horizon.
_Avoid_: Filter, content policy

**Canonical Status**:
The declared relationship of a Claim or Work to a selected continuity.
_Avoid_: Truth level, officialness

**Continuity**:
A coherent interpretation of Works and Claims that may coexist without an
assumed contradiction.
_Avoid_: Universe, canon

