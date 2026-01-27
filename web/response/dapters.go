package response

func (r *Response) GJSON(c interface{ JSON(int, interface{}) }) {
	c.JSON(r.Status(), r)
}
