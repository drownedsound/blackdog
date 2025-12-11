-- TODO: Rename application_number to member_reference_no

.mode quote

SELECT 
    json_object(
    'application_number', app.application_number,
    'status', app.status_code,
    'requested_amount', app.requested_amount,
    'created_at', app.created_at,

    -- Product Details (Card or Loan)
    'product_details', json_object(
    'card_design', cd.card_design_code,
    'credit_limit', cd.credit_limit,
    'loan_term', ld.term_months,
    'loan_rate', ld.interest_rate
    ),

    -- Nested List of Applicants
    'applicants', (
        SELECT json_group_array(
            json_object(
            'role', apl.role_code,
            'product', apl.product_type_code,
            'first_name', apl.first_name,
            'last_name', apl.last_name,
            'dob', apl.date_of_birth,

            'identifications', (
                SELECT json_group_array(json_object(
                    'type', idt.id_type_code,
                    'value', idt.id_value,
                    'authority', idt.issuing_authority
                    ))
                FROM APPLICANT_IDENTIFICATION idt
                WHERE idt.applicant_id = apl.id
            ),

            'addresses', (
                SELECT json_group_array(json_object(
                    'type', adr.address_type_code,
                    'line_1', adr.line_1,
                    'line_2', adr.line_2,
                    'city', adr.city_municipality,
                    'province', adr.province,
                    'zip', adr.postal_code
                    ))
                FROM APPLICANT_ADDRESS adr
                WHERE adr.applicant_id = apl.id
            ),

            'contacts', (
                SELECT json_group_array(json_object(
                    'type', ctc.contact_type_code,
                    'value', ctc.contact_number
                    ))
                FROM APPLICANT_CONTACT ctc
                WHERE ctc.applicant_id = apl.id
            ),

            'employment', (
                SELECT json_group_array(json_object(
                    'status', emp.employment_status_code,
                    'employer', emp.employer_name,
                    'income_cents', emp.gross_monthly_income
                    ))
                FROM APPLICANT_EMPLOYMENT emp
                WHERE emp.applicant_id = apl.id
            )
            )
            )
        FROM APPLICANT apl
        WHERE apl.application_id = app.id
    )
    ) as application_json
FROM APPLICATION app
LEFT JOIN CARD_DETAIL cd ON app.id = cd.application_id
LEFT JOIN LOAN_DETAIL ld ON app.id = ld.application_id
WHERE app.application_number IN ('APP-2023-SINGLE', 'APP-2023-MULTI');
