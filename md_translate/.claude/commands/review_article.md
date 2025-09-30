# Review Article

You are tasked with conducting comprehensive review and refinement of article content, ensuring quality, accuracy, engagement, and alignment with objectives before final publication.

## CRITICAL: YOUR ROLE IS TO REVIEW AND REFINE EXISTING ARTICLE CONTENT
- DO NOT rewrite the entire article structure
- DO NOT conduct new research (use existing materials for fact-checking)
- DO NOT change the core message or planned content strategy
- FOCUS on improving clarity, accuracy, flow, and engagement
- You are creating the polished, publication-ready version

## Initial Setup:

When this command is invoked, respond with:
```
I'm ready to review your article. Please provide:
1. Path to the article draft (from write_article phase)
2. Path to plan document (for reference)
3. Path to research document (for fact-checking)
4. Specific review focus areas (content, style, accuracy, engagement, etc.)
5. Target publication standards or requirements
6. Any style guide or formatting requirements
```

Then wait for the user's review requirements.

## Steps to follow after receiving the review requirements:

1. **Read article draft and supporting documents:**
   - Read the article draft FULLY first
   - **IMPORTANT**: Use the Read tool WITHOUT limit/offset parameters
   - Read plan document to understand original objectives
   - Read research document for fact-checking reference
   - **CRITICAL**: Read all referenced materials before beginning review

2. **Analyze article and create review strategy:**
   - Assess overall structure and flow
   - Evaluate content quality and accuracy
   - Check alignment with original plan and objectives
   - Identify areas needing improvement
   - Create review plan using TodoWrite to track review tasks

3. **Conduct comprehensive content review:**
   - Perform systematic review across multiple dimensions:
   
   **Content Quality Review:**
   - Verify accuracy of all facts, statistics, and claims
   - Check proper integration of research findings
   - Ensure expert quotes are accurate and contextual
   - Validate that all major points are well-supported
   
   **Structure and Flow Review:**
   - Assess logical progression between sections
   - Evaluate transition quality and smoothness
   - Check introduction effectiveness and conclusion impact
   - Ensure each section fulfills its planned purpose
   
   **Audience and Engagement Review:**
   - Verify content matches target audience level
   - Assess tone and style consistency
   - Evaluate engagement elements (hooks, examples, stories)
   - Check clarity of complex concepts and explanations
   
   **Technical Review:**
   - Check formatting, headings, and structure
   - Verify source citations and references
   - Ensure proper grammar, spelling, and style
   - Validate word count against targets

4. **Create detailed review report:**
   - Document all findings and recommended improvements
   - Prioritize changes by impact and importance
   - Provide specific suggestions with rationale
   - Include both major structural issues and minor edits

5. **Implement agreed-upon improvements:**
   - Make revisions based on review findings
   - Focus on high-impact improvements first
   - Maintain original voice and style while enhancing clarity
   - Ensure all changes align with article objectives

6. **Gather metadata for final article:**
   - Run the `hack/spec_metadata.sh` script to generate relevant metadata
   - Update filename to reflect final status:
     - `articles/YYYY-MM-DD-topic-final.md` for completed articles
     - Or maintain user-specified location with version update

7. **Generate review report and final article:**
   ```markdown
   ---
   date: [Current date and time with timezone in ISO format]
   reviewer: [Reviewer name from thoughts status]
   git_commit: [Current commit hash]
   branch: [Current branch name]
   repository: [Repository name]
   title: "[Article Title]"
   topic: "[Article Topic]"
   target_audience: "[Target audience]"
   article_format: "[Format type]"
   word_count: [Final word count]
   draft_document: "[Path to original draft]"
   plan_document: "[Path to plan document]"
   research_document: "[Path to research document]"
   tags: [article, final, reviewed, topic-specific-tags]
   status: reviewed
   last_updated: [Current date in YYYY-MM-DD format]
   last_updated_by: [Reviewer name]
   phase: review
   review_date: [Current date]
   ---

   # [Article Title]

   [Reviewed and refined article content]

   ---

   ## Review Report

   ### Review Summary
   **Review Date**: [Current date]
   **Original Word Count**: [Original count]
   **Final Word Count**: [Final count]
   **Review Focus**: [Areas of focus during review]

   ### Changes Made

   #### Major Improvements
   - [Significant change 1]: [Rationale and impact]
   - [Significant change 2]: [Rationale and impact]

   #### Content Enhancements
   - [Content improvement 1]: [Description]
   - [Content improvement 2]: [Description]

   #### Style and Flow Refinements
   - [Style improvement 1]: [Description]
   - [Style improvement 2]: [Description]

   #### Fact-Checking Results
   - [Verification 1]: [Result]
   - [Verification 2]: [Result]

   ### Quality Metrics

   #### Content Quality
   - [ ] All facts verified against research sources
   - [ ] Expert quotes accurate and properly attributed
   - [ ] Claims supported by evidence
   - [ ] Current and relevant information throughout

   #### Structure and Flow
   - [ ] Logical progression from introduction to conclusion
   - [ ] Smooth transitions between sections
   - [ ] Each section fulfills its planned purpose
   - [ ] Compelling introduction and strong conclusion

   #### Audience Alignment
   - [ ] Appropriate complexity for target audience
   - [ ] Consistent tone and style throughout
   - [ ] Clear explanations of technical concepts
   - [ ] Engaging and accessible language

   #### Technical Standards
   - [ ] Proper formatting and structure
   - [ ] Accurate citations and references
   - [ ] Grammar and spelling checked
   - [ ] Word count within target range

   ### Publication Readiness
   - **Content Accuracy**: [Rating/Notes]
   - **Engagement Level**: [Rating/Notes]
   - **Target Audience Fit**: [Rating/Notes]
   - **Technical Quality**: [Rating/Notes]
   - **Overall Assessment**: [Ready for publication/Needs minor edits/Major revision needed]

   ### Recommendations for Publication
   1. [Recommendation 1]
   2. [Recommendation 2]
   3. [Any final suggestions]

   ### Sources Verification
   [Updated source list with verification status]
   - [Source 1]: Verified [date] - [URL/Reference]
   - [Source 2]: Verified [date] - [URL/Reference]
   ```

8. **Final quality assurance:**
   - Read through entire refined article once more
   - Check that all review items are addressed
   - Ensure consistency in style, tone, and formatting
   - Verify all sources and citations are accurate
   - Confirm article meets publication standards

9. **Sync and present final results:**
   - Run `humanlayer thoughts sync` if saving to thoughts directory
   - Present summary of review process and final article status
   - Highlight key improvements made during review
   - Provide publication readiness assessment
   - Offer final recommendations or next steps

## Review Focus Areas:

### Content Accuracy
- Fact-check all statistics, dates, and claims
- Verify expert quotes and attributions
- Cross-reference with research sources
- Ensure currency of information
- Check for any unsupported assertions

### Clarity and Readability
- Assess sentence structure and complexity
- Check for jargon or overly technical language
- Ensure smooth paragraph transitions
- Evaluate explanation quality for complex topics
- Test logical flow of arguments

### Engagement and Impact
- Review hook effectiveness in introduction
- Assess use of examples, stories, and case studies
- Evaluate call-to-action strength
- Check for engaging language and varied sentence structure
- Ensure value proposition is clear throughout

### Style and Consistency
- Maintain consistent tone throughout
- Check formatting and structural elements
- Verify adherence to style guidelines
- Ensure appropriate voice for target audience
- Review heading structure and organization

### Technical Quality
- Check grammar, spelling, and punctuation
- Verify proper citation format
- Ensure consistent formatting
- Check links and references
- Validate metadata accuracy

## Important notes:
- Balance thorough review with preserving original voice
- Focus on improvements that enhance reader value
- Maintain alignment with original article objectives
- Prioritize changes that improve clarity and engagement
- Ensure all fact-checking is thorough and accurate
- Keep target audience needs central to all review decisions
- Document rationale for significant changes made
- Create final version that meets publication standards
- Provide actionable feedback for any remaining improvements needed
- Ensure review process adds measurable value to the article